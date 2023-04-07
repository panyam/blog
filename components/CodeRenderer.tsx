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

class CodeRenderer extends React.Component {
  codeRef = React.createRef<HTMLElement>()
  state: any

  constructor(props, context) {
    super(props, context)
    this.state = {
      url: null,
      lineStart: 1,
      title: '',
      height: '300px',
      language: '',
      contents: '',
      ...props,
    }
  }

  componentDidMount() {
    const state = this.state
    const language = state.language
    console.log('In Comp Did Mount: ', state.url)
    if (state.url) {
      console.log('Making reuqest: ', state.url)
      axios
        .get(state.url)
        .then(async (response) => {
          console.log('Got Resonse: ', response.data)
          this.setState((ps) => ({
            ...ps,
            contents: response.data,
          }))

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
  }

  render() {
    const state = this.state
    return (
      <div className="CodeEmbedContainer">
        <h4 className="CodeEmbedHeading">
          <Link className="CodeEmbedUrlLink" href={state.url}>
            {state.title}
          </Link>
        </h4>
        <div style={{ overflow: 'scroll', maxHeight: state.height || '350px' }}>
          <pre className="line-numbers" data-start={state.lineStart || 1}>
            <code ref={this.codeRef} className={`CodeEmbedCodeElem language-${state.language}`}>
              {`${state.contents}`}
            </code>
          </pre>
        </div>
      </div>
    )
  }
}

const CodeRenderer2 = (props) => {
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
