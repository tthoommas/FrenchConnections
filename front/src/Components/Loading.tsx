interface LoadingProps {
    message?: string
}

export default function Loading({ message = "Chargement..." }: LoadingProps) {
    return <div className="d-flex flex-column align-items-center justify-content-center vw-100 vh-100">
        <div className="spinner-grow text-primary" role="status">
            <span className="visually-hidden">{message}</span>
        </div>
        <p className="m-3 fw-bold fs-1">{message}</p>
    </div>
}