interface WordProps {
    word: string
    onClick: () => void
    selected: boolean
    editable?: boolean
    onWordEdited?: (newWord: string) => void
}

export default function Word({ word, onClick, selected, editable = false, onWordEdited = () => {} }: WordProps) {
    return <div className={`ratio ratio-1x1 border border-2 rounded-3 fw-medium ${selected ? "bg-dark-subtle" : "bg-light"}`}
        onClick={onClick}>
        <div className="d-flex align-items-center justify-content-center">
            {
                editable ? <input type="text"
                className={`form-control text-center bg-transparent border-0 no-focus-ring fw-medium text-uppercase`}
                placeholder="Mot"
                value={word}
                readOnly={!editable}
                style={{ fontSize: "0.9em" }} 
                onChange={(e) => onWordEdited(e.target.value)}/>:
                <span className="text-center text-truncate text-uppercase fw-medium user-select-none">{word}</span>
            }
            
        </div>
    </div>
}