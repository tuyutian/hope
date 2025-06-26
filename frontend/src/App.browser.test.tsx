import App from "./App";
import {expect, test} from "vitest";
import {render} from "vitest-browser-react";

test("renders learn react link", async () => {

  const {getByText} = render(<App />);
  await expect.element(getByText(/learn react/i)).toBeInTheDocument();

});
