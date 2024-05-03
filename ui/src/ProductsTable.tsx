import { Table, TableContainer, Tbody, Th, Thead, Tr } from "@chakra-ui/table";
import { Product } from "product";
import ProductRow from "ProductRow";
import { FunctionComponent, useEffect, useState } from "react";

export interface ProductsTableProps {
  products: Product[];
  handleDelete: (id: number) => void;
}

const ProductsTable: FunctionComponent<ProductsTableProps> = ({
  products: initialProducts,
  handleDelete,
}) => {
  const [products, setProducts] = useState<Product[]>(initialProducts);

  useEffect(() => {
    setProducts(initialProducts);
  }, [initialProducts]);

  return (
    <TableContainer>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th>Id</Th>
            <Th>Brand</Th>
            <Th>Category</Th>
            <Th>Quantity</Th>
            <Th>Price</Th>
          </Tr>
        </Thead>
        <Tbody>
          {products.map((item, index) => (
            <ProductRow
              key={index}
              product={item}
              onDelete={() => handleDelete(item.id)}
            />
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default ProductsTable;
