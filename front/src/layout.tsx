import { Outlet, useNavigate } from "react-router";


export default function Layout() {
    const navigate = useNavigate()
    return <>
        <main>
            <nav className="navbar navbar-expand-lg bg-body-tertiary">
                <div className="container-fluid">
                    <button className="border-0 bg-transparent navbar-brand">FrenchConnections</button>
                    <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>
                    <div className="collapse navbar-collapse" id="navbarSupportedContent">
                        <ul className="navbar-nav me-auto mb-2 mb-lg-0">
                            <li className="nav-item">
                                <button className="nav-link active" aria-current="page" onClick={() => navigate("/")}>Jouer</button>
                            </li>
                            <li className="nav-item">
                                <button className="nav-link" onClick={() => navigate("/create")}>Créer</button>
                            </li>
                        </ul>
                    </div>
                </div>
            </nav>
            <Outlet />
        </main>

        <footer className="w-100 mt-5">
            <div className="container-fluid">
                <p className="m-0 p-0">Original <a target="_blank" rel="noreferrer" href="https://www.nytimes.com/games/connections">game</a></p>
                <p className="m-0 p-0">Opensource <a target="_blank" rel="noreferrer" href="https://github.com/tthoommas/FrenchConnections">Github project</a></p>
            </div>
        </footer>
    </>
}