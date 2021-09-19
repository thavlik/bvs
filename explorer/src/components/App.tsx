import * as React from "react";
import { hot } from "react-hot-loader";

import "./../assets/scss/App.scss";

class App extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
      <div className="app">
        <h1>Search Results</h1>
        <p>TODO: mock up search results</p>
      </div>
    );
  }
}

declare let module: Record<string, unknown>;

export default hot(module)(App);
