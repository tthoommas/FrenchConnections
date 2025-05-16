import { useEffect, useState } from "react"
import ShuffledGame from "../models/shuffled_game"
import { useNavigate } from "react-router"

export default function List() {

    const navigate = useNavigate()
    const [games, setGames] = useState<ShuffledGame[]>([])
    useEffect(() => {
        fetch(process.env.REACT_APP_API_ROOT + `/game/list`).then((resp) => {
            if (resp.ok) {
                return resp.json()
            }
        }).then((games: ShuffledGame[]) => {
            setGames(games)
        })
    }, [])
    return <div className="container mt-4">
        {
            games.map((game) => {
                const date = new Date(game.createdAt);
                return <div className="row border border-2 rounded-3 fw-medium bg-light p-2 m-2 clickable" onClick={() => navigate(`/play/${game.id}`)}> 
                    Jeu n° {game.id} créé par {game.createdBy} le {date.toLocaleDateString()}
                </div>
            })
        }
    </div>
}