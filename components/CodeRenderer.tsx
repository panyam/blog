import React from 'react'
import Link from './Link'
import axios from 'axios'
import 'prismjs/plugins/line-numbers/prism-line-numbers.css'

/*
async function highlight(language: string, codeRef: React.RefObject<HTMLElement>) {
  if (typeof window !== 'undefined' || !language) {
    //import the language dynamically using import statement
    await import(`prismjs/components/prism-${language}`)
    Prism.highlightElement(codeRef.current)
  }
}
*/

const CodeRenderer = (props) => {
  const language = props.language
  const [contents, setContents] = React.useState(props.contents || '')
  const codeRef = React.createRef<HTMLElement>()

  React.useEffect(() => {
    if (props.url) {
      axios
        .get(props.url)
        .then(async (response) => {
          setContents(response.data.trim())
          if (typeof window !== 'undefined' || !language) {
            const Prism = await import('prismjs')
            await import('prismjs/plugins/line-numbers/prism-line-numbers.js')
            try {
              await import(`prismjs/components/prism-${language}`)
            } catch (err) {
              console.log('Error loading lang plugin: ', language)
            }
            // const t1 = Date.now()
            // Prism.highlightElement(codeRef.current)
            Prism.highlightAll()
            // console.log("Highlighting: ", Date.now() - t1, codeRef.current)
          }
        })
        .catch((err) => {
          console.log(err)
        })
    }
  }, [contents, language])

  return (
    <div>
      <h4>
        <Link href={props.url}>{props.title}</Link>
      </h4>
      <div style={{ overflow: 'scroll', maxHeight: props.height || '350px' }}>
        <pre className="line-numbers" data-start={props.lineStart || 1}>
          <code ref={codeRef} className={`language-${language}`}>{`${contents}`}</code>
        </pre>
      </div>
    </div>
  )
}

export default CodeRenderer
