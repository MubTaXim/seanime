import { MangaChapterDetails, MangaEntry } from "@/app/(main)/manga/_lib/manga.types"
import { __manga_selectedChapterAtom } from "@/app/(main)/manga/entry/_containers/chapter-reader/chapter-reader-drawer"
import { serverStatusAtom } from "@/atoms/server-status"
import { SeaEndpoints } from "@/lib/server/endpoints"
import { useSeaMutation } from "@/lib/server/query"
import { useAtomValue } from "jotai/react"
import React from "react"

export type DiscordPresenceRoute_QueryVariables = {
    title: string
    image: string
    chapter: string
}

export function useDiscordMangaPresence(entry: MangaEntry | undefined) {
    const serverStatus = useAtomValue(serverStatusAtom)
    const currentChapter = useAtomValue(__manga_selectedChapterAtom)

    const { mutate } = useSeaMutation<boolean, DiscordPresenceRoute_QueryVariables>({
        endpoint: SeaEndpoints.DISCORD_PRESENCE_MANGA,
    })
    const { mutate: cancelActivity } = useSeaMutation({
        endpoint: SeaEndpoints.DISCORD_PRESENCE_CANCEL,
    })

    const prevChapter = React.useRef<MangaChapterDetails | undefined>()

    React.useEffect(() => {
        if (
            serverStatus?.settings?.discord?.enableRichPresence &&
            serverStatus?.settings?.discord?.enableMangaRichPresence
        ) {

            if (currentChapter && entry) {
                mutate({
                    title: entry.media.title?.userPreferred || entry.media.title?.romaji || entry.media.title?.english || "Reading",
                    image: entry.media.coverImage?.large || entry.media.coverImage?.medium || "",
                    chapter: currentChapter.chapter,
                })
            }

            if (!currentChapter) {
                cancelActivity()
            }
        }

        prevChapter.current = currentChapter
    }, [currentChapter, entry])
}