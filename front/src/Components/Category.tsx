import Category from "../models/category"

interface CategoryProps {
    category: Category
}

export default function DisplayCategory({ category }: CategoryProps) {
    return <div className="d-flex flex-column border bg-success rounded text-white align-items-center justify-content-center m-1">
        <p className="fw-bold m-0 p-2">{category.categoryTitle}</p>
        <p className="m-0 p-1">{category.words.join(",")}</p>
    </div>
}