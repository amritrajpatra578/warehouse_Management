import {
  Button,
  FormControl,
  FormLabel,
  Input,
  useToast,
} from "@chakra-ui/react";
import ErrorList from "ErrorList";
import router from "next/router";
import { Product } from "product";
import { ChangeEvent, FunctionComponent, useState } from "react";

export interface ProductFormProps {
  onSubmit: (product: Product) => Promise<void>;
  initialProduct?: Product;
}

function handleProduct(product?: Product) {
  if (product != undefined) {
    return product;
  }
  return {
    id: 0,
    brand: "",
    category: "",
    price: 0,
    quantity: 0,
  };
}
const ProductForm: FunctionComponent<ProductFormProps> = ({
  onSubmit,
  initialProduct,
}) => {
  const [product, setProduct] = useState<Product>(
    handleProduct(initialProduct)
  );
  const toast = useToast();

  const handleSubmit = () => {
    onSubmit(product)
      .then(() => {
        toast({
          title: "saved",
          isClosable: true,
          status: "success",
        });
        router.push("http://localhost:8000/listproducts");
      })
      .catch((error) => {
        const validationErrors = error.response?.data?.errors;
        toast({
          title: error.message,
          description: <ErrorList errors={validationErrors} />,
          isClosable: true,
          status: "error",
        });
      });
  };

  return (
    <FormControl>
      <FormLabel>ID :</FormLabel>
      <Input
        disabled={initialProduct != undefined}
        type="number"
        name="id"
        value={product?.id}
        onChange={(e: ChangeEvent<HTMLInputElement>) => {
          setProduct({
            ...product,
            id: parseInt(e.target.value),
          });
        }}
      />
      <FormLabel>Brand :</FormLabel>
      <Input
        type="text"
        value={product?.brand}
        onChange={(e: ChangeEvent<HTMLInputElement>) => {
          setProduct({ ...product, brand: e.target.value });
        }}
      />
      <FormLabel>Category :</FormLabel>
      <Input
        type="text"
        value={product?.category}
        onChange={(e: ChangeEvent<HTMLInputElement>) => {
          setProduct({ ...product, category: e.target.value });
        }}
      />
      <FormLabel>Quantity :</FormLabel>
      <Input
        type="number"
        value={product?.quantity}
        onChange={(e: ChangeEvent<HTMLInputElement>) => {
          setProduct({
            ...product,
            quantity: parseInt(e.target.value),
          });
        }}
      />
      <FormLabel>Price :</FormLabel>
      <Input
        type="number"
        value={product?.price}
        onChange={(e: ChangeEvent<HTMLInputElement>) => {
          setProduct({
            ...product,
            price: parseInt(e.target.value),
          });
        }}
      />
      <Button colorScheme="green" onClick={handleSubmit}>
        Save
      </Button>
    </FormControl>
  );
};

export default ProductForm;
