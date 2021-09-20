import * as React from "react";

import { BigSearchBar } from "./BigSearchBar";

import "./HomePage.scss";

export class HomePage extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
      <div className="home-page">
          <h1>The future of election security is here.</h1>
          <BigSearchBar />
          TODO: show links, stats, and graphs for any ongoing elections
      </div>
    );
  }
}

