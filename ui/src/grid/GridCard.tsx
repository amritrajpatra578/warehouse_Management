import { FunctionComponent } from "react";
import Grid from "./Grid";
import { GridItem } from "./grid";
import styles from "styles/Home.module.css";
import { Inter } from "next/font/google";

const inter = Inter({ subsets: ["latin"] });

export interface GridCardProps {
  gridItem: GridItem;
}

const GridCard: FunctionComponent<GridCardProps> = ({ gridItem }) => {
  return (
    <a
      href={gridItem.href}
      className={styles.card}
      target="_blank"
      rel="noopener noreferrer"
    >
      <h2 className={inter.className}>
        {gridItem.title} <span>-&gt;</span>
      </h2>
      <p className={inter.className}>{gridItem.description}</p>
    </a>
  );
};

export default GridCard;
