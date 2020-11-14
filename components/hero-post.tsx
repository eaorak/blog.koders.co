import Avatar from "./avatar"
import DateFormater from "./date-formater"
import CoverImage from "./cover-image"
import Link from "next/link"
import Author from "../types/author"

type Props = {
  title: string
  coverImage: string
  date: string
  excerpt: string
  author: Author
  slug: string
}

const HeroPost = ({ title, coverImage, date, excerpt, author, slug }: Props) => {
  return (
    <section>
      <div className='mb-8 md:mb-16'>
        <CoverImage title={title} src={coverImage} slug={slug} />
      </div>
      <div>
        <h3 className='mb-4 text-4xl lg:text-6xl leading-tight'>
          <Link as={`/posts/${slug}`} href='/posts/[slug]'>
            <a className='hover:underline'>{title}</a>
          </Link>
        </h3>
        <div className='mb-4 text-lg'>
          <DateFormater dateString={date} />
        </div>
      </div>
      <div className='mb-16'>
        <p className='text-lg leading-relaxed'>{excerpt}</p>
        <Avatar name={author.name} picture={author.picture} />
      </div>
    </section>
  )
}

export default HeroPost
