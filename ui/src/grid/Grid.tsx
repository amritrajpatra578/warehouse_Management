import { FunctionComponent } from "react";
import { GridItem } from "./grid";
import GridCard from "./GridCard";
import styles from "styles/Home.module.css";
import { Inter } from "next/font/google";
const inter = Inter({ subsets: ["latin"] });

export interface GridProps {
  items: GridItem[];
}

const Grid: FunctionComponent<GridProps> = ({ items }) => {
  return (
    <div className={styles.grid}>
      {items.map((item, index) => (
        <GridCard key={index} gridItem={item} />
      ))}
    </div>
  );
};

export default Grid;
