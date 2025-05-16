import { useState } from "react";
import Category from "../models/category";
import Word from "../Components/Word";
import CreateResponse from "../models/createResponse";
import { useNavigate } from "react-router";

export default function Create() {
    const navigate = useNavigate()
    const [createdBy, setCreatedBy] = useState("")
    const [gameCategories, setGamecategories] = useState<Category[]>([
        { categoryTitle: "", words: ["", "", "", ""] },
        { categoryTitle: "", words: ["", "", "", ""] },
        { categoryTitle: "", words: ["", "", "", ""] },
        { categoryTitle: "", words: ["", "", "", ""] }
    ])
    

    const [errorMessage, setErrorMessage] = useState("")

    const onCategoryEdited = (catIndex: number, newCategory: string) => {
        setGamecategories((gameCategories) => {
            let copy = [...gameCategories]
            copy[catIndex].categoryTitle = newCategory
            return copy
        })
    }

    const onWordEdited = (catIndex: number, wordIndex: number, newWord: string) => {
        console.log(catIndex, wordIndex, newWord)
        setGamecategories((gameCategories) => {
            let copy = [...gameCategories]
            copy[catIndex].words[wordIndex] = newWord
            return copy
        })
    }

    const onSubmitNewGame = () => {
        let seenCategories = new Set()
        let seenWords = new Set()
        for (let category of gameCategories) {
            if (category.categoryTitle.trim().length === 0) {
                setErrorMessage("Au moins un titre de catégorie est manquant")
                return
            }
            if (seenCategories.has(category.categoryTitle)) {
                setErrorMessage("Les catégories doivent être uniques")
                return
            } else {
                seenCategories.add(category.categoryTitle)
            }
            for (let word of category.words) {
                if (word.trim().length === 0) {
                    setErrorMessage("Au moins un mot est manquant")
                    return
                }
                if (seenWords.has(word)) {
                    setErrorMessage("Un même mot ne doit apparaître qu'une seule fois dans l'ensemble du jeu")
                    return
                } else {
                    seenWords.add(word)
                }
            }
        }
        if (createdBy.trim().length < 1) {
            setErrorMessage("Renseignez un nom")
            return
        }

        setErrorMessage("")

        fetch(process.env.REACT_APP_API_ROOT + `/game`, {
            method: "POST",
            body: JSON.stringify({ createdBy, gameCategories})
        }).then((resp) => {
            if (resp.ok) {
                return resp.json()
            }
        }).then((result: CreateResponse) => {
            navigate(`/play/${result.gameId}`)
        }).catch((err) => {
            setErrorMessage("Erreur le jeu n'a pas pu être créé")
        })
    }

    return <div className="container mt-3">
        <div className="row">
            <div className="col text-center">
                <h1 className="fs-2 my-3">Nouveau jeu</h1>
            </div>
        </div>
        <div className="row">
            <div className="container grid mt-3">
                <div className="row row-cols-4 g-2 g-sm-2">
                    {
                        gameCategories.map((cat, catIndex) => {
                            return <>
                                <div className="col-12 border border-2 rounded-3 user-select-none fw-medium bg-light">
                                    <input
                                        type="text"
                                        className="form-control bg-transparent border-0 no-focus-ring"
                                        placeholder="Nom de la catégorie"
                                        value={gameCategories[catIndex].categoryTitle}
                                        onChange={(e) => onCategoryEdited(catIndex, e.target.value)}
                                    />
                                </div>
                                {
                                    cat.words.map((word, wordIndex) => {
                                        return <div className="col" key={wordIndex}>
                                            <Word
                                                word={word}
                                                onClick={() => { }}
                                                selected={false}
                                                editable={true}
                                                onWordEdited={(newWord) => onWordEdited(catIndex, wordIndex, newWord)}
                                            />
                                        </div>
                                    })
                                }
                            </>
                        })

                    }
                    <div className="col-12 border border-2 rounded-3 user-select-none fw-medium bg-light mt-5">
                        <input
                            type="text"
                            className="form-control bg-transparent border-0 no-focus-ring"
                            placeholder="Créé par ..."
                            value={createdBy}
                            onChange={(e) => setCreatedBy(e.target.value)}
                        />
                    </div>
                </div>
            </div>
        </div>

        {
            errorMessage.length > 0 && <div className="row">
                <div className="alert alert-warning my-3" role="alert">
                    {errorMessage}
                </div>
            </div>
        }
        <div className="row">
            <div className="col text-center">
                <button type="button" onClick={onSubmitNewGame} className="btn btn-primary btn m-3">Valider</button>
            </div>
        </div>
    </div>
}