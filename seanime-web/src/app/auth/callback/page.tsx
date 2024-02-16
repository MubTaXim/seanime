"use client"
import { useAuth } from "@/app/auth/_lib/auth"
import { serverStatusAtom } from "@/atoms/server-status"
import { LoadingOverlay } from "@/components/ui/loading-spinner"
import { useSetAtom } from "jotai/react"
import { useRouter } from "next/navigation"
import React, { useEffect } from "react"
import toast from "react-hot-toast"
import { useUpdateEffect } from "react-use"

export default function Page() {

    const router = useRouter()

    const { data, error, token } = useAuth()
    const setServerStatus = useSetAtom(serverStatusAtom)

    useUpdateEffect(() => {
        if (error) {
            toast.error(error.message)
            router.push("/")
        }
    }, [error])

    useEffect(() => {
        if (window !== undefined && !!data && !!token) {
            setServerStatus(data)
            const t = setTimeout(() => {
                toast.success("Successfully authenticated")
                router.push("/")
            }, 1000)

            return () => clearTimeout(t)
        }
    }, [data])

    return (
        <div>
            <LoadingOverlay className={"fixed w-full h-full z-[80]"}>
                <h3 className={"mt-2"}>Authenticating...</h3>
            </LoadingOverlay>
        </div>
    )
}
