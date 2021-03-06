import { CMS_NAME } from "../lib/constants"

const Intro = () => {
  return (
    <section className='flex-col md:flex-row flex items-center md:justify-between mt-16 mb-16 md:mb-12'>
      <h1 className='text-4xl md:text-4xl font-bold tracking-tighter leading-tight md:pr-8'>
        <a href="https://koders.co"><img src='/assets/images/koders.png' alt='koders' style={{ display: "inline-block", maxHeight: "120px" }} /></a>
      </h1>
      <h4 className='text-center md:text-left text-lg mt-5 md:pl-8'>Follow the Tech world with us.</h4>
    </section>
  )
}

export default Intro
