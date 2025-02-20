import { render, screen } from "@testing-library/react";
import App from "./App";

test("renders react text", () => {
  render(<App />);
  const linkElement = screen.getByText(/react page/i);
  expect(linkElement).toBeInTheDocument();
});

test("renders velocity text", () => {
  render(<App />);
  const linkElement = screen.getByText(/learn velocity/i);
  expect(linkElement).toBeInTheDocument();
});
