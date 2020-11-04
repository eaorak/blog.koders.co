import Link from "next/link"

const Header = () => {
  return (
    <h2 className='text-2xl md:text-4xl font-bold tracking-tight md:tracking-tighter leading-tight mb-8 mt-8'>
      <Link href='/'>
        <img
          src='/assets/images/m.png'
          alt='meftun'
          style={{ display: "inline-block", maxHeight: 50, margin: 0, padding: 0 }}
          className='cursor-pointer'
        />
      </Link>
      .
    </h2>
  )
}

export default Header
