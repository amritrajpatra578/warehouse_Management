import { Button } from "@chakra-ui/react";
import { Td, Tr } from "@chakra-ui/table";
import Link from "next/link";
import { Product } from "product";
import React, { FunctionComponent } from "react";

export interface ProductRowProps {
  product: Product;
  onDelete: () => void;
}

const ProductRow: FunctionComponent<ProductRowProps> = ({
  product,
  onDelete,
}) => {
  return (
    <>
      <Tr>
        <Td>{product.id}</Td>
        <Td>{product.brand}</Td>
        <Td>{product.category}</Td>
        <Td>{product.quantity}</Td>
        <Td>{product.price}</Td>
        <Td>
          <Link
            href={{
              pathname: "http://localhost:8000/updateproduct/[slug]",
              query: { slug: product.id },
            }}
          >
            <Button colorScheme="teal" variant="outline" cursor="pointer">
              Update
            </Button>
          </Link>
        </Td>
        <Td>
          <Button
            colorScheme="red"
            variant="outline"
            onClick={() => {
              onDelete();
            }}
          >
            Delete
          </Button>
        </Td>
      </Tr>
    </>
  );
};

export default ProductRow;
