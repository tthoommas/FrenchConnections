import Word from "./Word"

interface GridProps {
    words: string[]
    onWordClicked: (word: string) => void
    selectedWords: Set<String>
    isWrongGuess?: boolean
    isOneAway?: boolean
}

export default function Grid({ words, onWordClicked, selectedWords, isWrongGuess = false, isOneAway = false }: GridProps) {
    return <div className="container grid mt-3">
        <div className="row row-cols-4 g-2 g-sm-2">
            {
                words.map((word) => {
                    let isSelected = selectedWords.has(word)
                    return <div className="col" key={word}>
                        <Word word={word} onClick={() => onWordClicked(word)} selected={isSelected} isOneAway={isOneAway} isWrongGuess={isWrongGuess} />
                    </div>
                })
            }
        </div>
    </div>
}