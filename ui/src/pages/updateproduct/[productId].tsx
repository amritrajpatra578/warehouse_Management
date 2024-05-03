import { Spinner, useToast } from "@chakra-ui/react";
import { getProductById, updateProduct } from "api";
import { useRouter } from "next/router";
import { Product } from "product";
import ProductForm from "ProductForm";
import { useEffect, useState } from "react";

interface UpdatePageState {
  product?: Product;
  loading: boolean;
  statusCode: number;
  error?: string;
}

export default function ProductDetails() {
  const router = useRouter();
  const [state, setState] = useState<UpdatePageState>({
    statusCode: 0,
    loading: true,
  });
  const toast = useToast();

  function fetchProduct() {
    const productId = router.query.productId as string;
    getProduct(parseInt(productId));
  }

  function getProduct(id: number) {
    getProductById(id)
      .then((val) => {
        setState({
          product: val,
          statusCode: 200,
          loading: false,
          error: undefined,
        });
      })
      .catch((error) => {
        console.log(error);
        setState({
          statusCode: 0,
          loading: false,
          error: error,
        });
        toast({
          title: error.message,
          isClosable: true,
          status: "error",
        });
      });
  }

  useEffect(() => {
    if (!router.isReady) {
      return;
    }
    fetchProduct();
  }, [router.isReady]);

  if (!router.isReady || state.loading) {
    return (
      <Spinner
        thickness="4px"
        speed="0.65s"
        emptyColor="gray.200"
        color="blue.500"
        size="xl"
      />
    );
  }

  if (state.error) {
    return state.error;
  }

  return (
    <ProductForm onSubmit={updateProduct} initialProduct={state.product} />
  );
}
