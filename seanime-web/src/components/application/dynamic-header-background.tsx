"use client"
import { cn } from "@/components/ui/core"
import Image from "next/image"
import { usePathname } from "next/navigation"
import React from "react"
import { SiAnilist } from "react-icons/si"

export function DynamicHeaderBackground() {

    const pathname = usePathname()

    return (
        <>
            {!pathname.startsWith("/entry") && <>
                {!pathname.startsWith("/anilist") && <div
                    className={cn(
                        "bg-[url(/pattern-2.svg)] bg-[--background-color] opacity-60 bg-cover bg-center bg-repeat z-[-2] w-full h-[10rem] absolute bottom-0",
                    )}
                />}
                {(pathname.startsWith("/anilist") && !pathname.startsWith("/search")) && <div
                    className={cn(
                        "bg-[url(/pattern-3.svg)] bg-blue-700/10 opacity-60 bg-contain bg-center bg-repeat z-[-2] w-full h-[20rem] absolute bottom-0",
                    )}
                >
                    <div className="w-full flex items-center justify-center absolute bottom-0 h-[10rem]">
                        <SiAnilist className="text-5xl text-white relative z-[1]" />
                    </div>
                </div>}
                {(pathname === "/") && <Image
                    src={"/landscape-tenki-no-ko.jpg"}
                    alt={"tenki no ko"}
                    fill
                    priority
                    className={"object-cover object-bottom opacity-30 z-[-2]"}
                />}
                {(pathname.startsWith("/search")) && <Image
                    src={"/landscape-tenki-no-ko.jpg"}
                    alt={"tenki no ko"}
                    fill
                    priority
                    className={"object-cover opacity-30 z-[-2]"}
                />}
                <div
                    className={"w-full absolute bottom-0 h-[10rem] bg-gradient-to-t from-[--background-color] to-transparent z-[-2]"}
                />
            </>}
        </>
    )
}
