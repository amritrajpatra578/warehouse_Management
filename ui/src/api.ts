import axios from "axios";
import { Product } from "product";

export async function getProducts(): Promise<Product[]> {
  const response = await axios.get("/products");
  return response.data;
}

export async function createProduct(product: Product): Promise<void> {
  await axios.post("/products", product);
}
export async function getProductById(id: number): Promise<Product> {
  const response = await axios.get(`/products/${id}`);
  return response.data;
}

export async function updateProduct(product: Product): Promise<void> {
  await axios.put(`/products/${product.id}`, product);
  console.log(product);
}

export async function deleteProduct(id: number): Promise<void> {
  await axios.delete(`/products/${id}`);
}
