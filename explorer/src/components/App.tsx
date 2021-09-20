import * as React from "react";
import { hot } from "react-hot-loader";
import { Banner } from "./Banner";
import "./../assets/scss/App.scss";

class App extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
      <div className="app">
        <Banner />
        <div className="content">
          <h1>The future of election security is here.</h1>
          <div className="main-search">
            <input type="text" value="Search for vote, election, minter..." />
            <div className="submit-search"></div>
          </div>
          TODO: show links, stats, and graphs for any ongoing elections
        </div>
        <div className="copyinfo">
          © 2021 Thomas Havlik. All rights reserved. <a href="https://github.com/thavlik/bvs">Open Source</a>
        </div>
      </div>
    );
  }
}

declare let module: Record<string, unknown>;

export default hot(module)(App);
