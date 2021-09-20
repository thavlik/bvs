import * as React from "react";
import { hot } from "react-hot-loader";
import { Banner } from "./Banner";
import "./../assets/scss/App.scss";
import { HomePage } from "./pages/home/HomePage";

class App extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
      <div className="app">
        <Banner />
        <div className="content">
          <HomePage />
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
