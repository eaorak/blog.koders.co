import Footer from "./footer"
import Meta from "./meta"
import { initAnalytics } from "./googleAnalytics.js"
import { useEffect } from "react"

type Props = {
  preview?: boolean
  children: React.ReactNode
}

const Layout = ({ children }: Props) => {
  useEffect(() => {
    initAnalytics()
  })

  return (
    <>
      <Meta />
      <div className='min-h-screen'>
        <main>{children}</main>
      </div>
      <Footer />
    </>
  )
}

export default Layout
