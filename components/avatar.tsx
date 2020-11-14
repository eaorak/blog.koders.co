type Props = {
  name: string
  picture: string
}

const Avatar = ({ name, picture }: Props) => {
  return (
    <div className='flex items-center'>
      <img src={picture} className='w-24 h-12 rounded-full mr-4 mt-4' alt={name} />
    </div>
  )
}

export default Avatar
