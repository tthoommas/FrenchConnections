import { useParams } from "react-router";
import Grid from "../Components/Grid";
import { useEffect, useState } from "react";
import ShuffledGame from "../models/shuffled_game";
import Loading from "../Components/Loading";
import GuessResponse from "../models/guessResponse";
import { eqSet } from "../utils";
import Category from "../models/category";
import DisplayCategory from "../Components/Category";

export default function GamePage() {
    let params = useParams()
    let [game, setGame] = useState<ShuffledGame | undefined>(undefined)
    const [selectedWords, setSelectedWords] = useState<Set<string>>(new Set())
    // Keep a set of 'one away' combinations already found
    const [oneAwaySets, setOneAwaySets] = useState<Set<string>[]>([])
    const [foundCategories, setFoundCategories] = useState<Category[]>([])
    const [tentatives, setTentatives] = useState(0)

    useEffect(() => {
        fetch(process.env.REACT_APP_API_ROOT + `/game/${params.gameId}`).then((resp) => {
            if (resp.ok) {
                return resp.json()
            }
        }).then((data) => {
            setOneAwaySets([])
            setFoundCategories([])
            setTentatives(0)
            setSelectedWords(new Set())
            setGame(data)
        })
    }, [params.gameId])

    const isAmongOneAway = (proposition: Set<string>, knownOneAways: Set<string>[]): boolean => {
        for (let oneAwaySet of knownOneAways) {
            if (eqSet(proposition, oneAwaySet)) {
                return true
            }
        }
        return false
    }

    const submitGuess = () => {
        if (selectedWords.size === 4) {
            fetch(process.env.REACT_APP_API_ROOT + `/game/${params.gameId}/guess`, {
                method: "POST",
                body: JSON.stringify({ "proposition": Array.from(selectedWords) })
            }).then((resp) => {
                if (resp.ok) {
                    return resp.json()
                }
            }).then((guessResponse: GuessResponse) => {
                setTentatives((old) => old + 1)
                if (guessResponse.isOneAway) {
                    if (!isAmongOneAway(selectedWords, oneAwaySets)) {
                        setOneAwaySets((oneAways) => [...oneAways, new Set(selectedWords)])
                    }
                } else if (guessResponse.success) {
                    setFoundCategories((alreadyFound) => {
                        return [...alreadyFound, { categoryTitle: guessResponse.categoryTitle, words: Array.from(selectedWords) }]
                    })
                    setSelectedWords(() => new Set())
                }
            })
        }
    }

    const onWordClicked = (word: string) => {
        setSelectedWords((selectedSet) => {
            // Toggle the word state
            const updatedSet = new Set(selectedSet)
            if (updatedSet.has(word)) {
                updatedSet.delete(word)
            } else if (updatedSet.size < 4) {
                updatedSet.add(word)
            }
            return updatedSet
        })
    }

    if (game === undefined) {
        return <Loading message="Chargement du jeu ..." />
    }

    let showOneAway = isAmongOneAway(selectedWords, oneAwaySets)
    let foundWords = foundCategories.flatMap(category => category.words)
    let done = foundCategories.length === 4

    return <div className="container">
        <div className="row">
            <h1 className="text-center fs-1 mt-3">Jeu n°{params.gameId}</h1>
            {game.createdBy && <p className="text-center">Créé par {game.createdBy}</p>}
        </div>
        <div className="row">
            <div className="col">
                {
                    foundCategories.map((cat) => {
                        return <DisplayCategory category={cat} key={cat.categoryTitle} />
                    })
                }
            </div>
        </div>
        <div className={`row my-3 ${!showOneAway ? "invisible" : ""}`}>
            <div className="col text-center">
                <span className="blink bg-dark-subtle p-1 border border-2 rounded-2">One away</span>
            </div>
        </div>
        <div className="row">
            <Grid words={game.game.filter((word) => !foundWords.includes(word))} onWordClicked={onWordClicked} selectedWords={selectedWords} />
        </div>
        {
            <div className="row">
                <div className="col text-center">
                    {
                        done ? <p className="fs-4">Terminé en {tentatives} tentatives !</p> : <button type="button" className="btn btn-secondary btn m-3" onClick={submitGuess}>Submit</button>
                    }
                </div>
            </div>
        }
    </div>
}