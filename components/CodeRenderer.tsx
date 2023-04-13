import React from 'react'
import Link from './Link'
import 'prismjs/plugins/line-numbers/prism-line-numbers.css'
import 'prismjs/plugins/line-numbers/prism-line-numbers.js'
import axios from 'axios'
import Prism from 'prismjs'

/*
async function highlight(language: string, codeRef: React.RefObject<HTMLElement>) {
  if (typeof window !== 'undefined' || !language) {
    //import the language dynamically using import statement
    await import(`prismjs/components/prism-${language}`)
    Prism.highlightElement(codeRef.current)
  }
}
*/

class CodeRendererClass extends React.Component {
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

  async componentDidMount() {
    const state = this.state
    const language = state.language
    console.log('In Comp Did Mount: ', state.url, language)
    if (state.url) {
      console.log('Making reuqest: ', state.url)
      axios
        .get(state.url)
        .then(async (response) => {
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

const CodeRendererFunc = async (props) => {
  const language = props.language
  // const [contents, setContents] = React.useState(props.contents || '')
  const codeRef = React.createRef<HTMLElement>()

  React.useEffect(() => {
    import(`prismjs/components/prism-${language}`)
      .then(() => {
        Prism.highlightAll()
      })
      .catch((err) => {
        console.log('Error loading lang plugin: ', language)
      })
  }, [props.url, language])

  console.log('Loading Code URL: ', props.url)
  let contents = props.contents || ''
  if (props.url) {
    const response = await axios.get(props.url)
    contents = response.data.trim()
    // const t1 = Date.now()
    // Prism.highlightElement(codeRef.current)
    // console.log("Highlighting: ", Date.now() - t1, codeRef.current)
    // if (typeof window !== 'undefined' || !language) { }
  }
  // React.useEffect(() => { }, [props.url, contents, language])

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

export default CodeRendererFunc
