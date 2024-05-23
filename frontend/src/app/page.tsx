import Navbar from "@/app/components/navbar";

const API_BASE_URL = "http://localhost:8080/api";
const API_VERSION = "v1";

export async function getProducts() {
    console.log(`${API_BASE_URL}/${API_VERSION}/products}`);
    const res = await fetch(`${API_BASE_URL}/${API_VERSION}/products`, {
        method: "GET",
    });

    if (!res.ok) {
        throw new Error(`${res.status}`);
    }
    return res.json();
}

interface Product {
    id: number,
    name: string,
    description: string,
    image: string,
    price: number,
    quantity: number,
    created_at: Date,
}

export default async function Page() {
    const products: Product[] = await getProducts();
    return (
        <main>
            <Navbar/>
            <div className="flex min-h-screen flex-col items-center p-24">
                <div>Welcome to shopping app</div>
                <div>{products.map((product) => product.name)} </div>
            </div>
        </main>
    );
}
