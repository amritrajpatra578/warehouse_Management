import { Spinner, useToast } from "@chakra-ui/react";
import { deleteProduct, getProducts } from "api";
import { Product } from "product";
import ProductsTable from "ProductsTable";
import { useEffect, useState } from "react";

export default function ProductsList() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const toast = useToast();

  function updateProductList() {
    getProducts()
      .then((products) => {
        setProducts(products);
      })
      .catch((error) => {
        toast({
          title: error.message,
          isClosable: true,
          status: "error",
        });
      })

      .finally(() => {
        setLoading(false);
      });
  }

  function handleDelete(id: number) {
    deleteProduct(id)
      .then(() => {
        updateProductList();
      })
      .catch((err) => {
        console.log(err);
      });
  }

  function connectWebsocket() {
    let socket = new WebSocket("ws://127.0.0.1:5000/ws");
    console.log("Attempting Connection...");

    socket.onopen = () => {
      console.log("Successfully Connected");
      socket.send("Hi From the Client!");
    };

    socket.onmessage = (message) => {
      console.log("data: ", message.data);
      var products = message.data;
      try {
        var result = JSON.parse(products);
        console.log(result);
        setProducts(result);
      } catch (err) {
        console.log("error while converting json to js object", err);
      }
    };

    socket.onclose = (event) => {
      console.log("Socket Closed Connection: ", event);
      socket.send("Client Closed!");
    };

    socket.onerror = (error) => {
      console.log("Socket Error: ", error);
    };
  }

  useEffect(() => {
    connectWebsocket();
    updateProductList();
  }, []);

  return (
    <>
      {loading && (
        <Spinner
          thickness="4px"
          speed="0.65s"
          emptyColor="gray.200"
          color="blue.500"
          size="xl"
        />
      )}
      <ProductsTable products={products} handleDelete={handleDelete} />
    </>
  );
}
