import { FunctionComponent } from "react";

export interface ErrorListProps {
  errors?: string[];
}

const ErrorList: FunctionComponent<ErrorListProps> = ({ errors }) => {
  if (!errors) {
    return null;
  }

  return (
    <ul>
      {errors.map((v: string, index: number) => (
        <li key={index}>{v}</li>
      ))}
    </ul>
  );
};

export default ErrorList;
