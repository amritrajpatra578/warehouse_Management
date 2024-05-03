import { createProduct } from "api";
import ProductForm from "ProductForm";

export default function CreateProductForm() {
  return <ProductForm onSubmit={createProduct} />;
}
