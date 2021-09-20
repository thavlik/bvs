import * as React from "react";

import "./BigSearchBar.scss";

export class BigSearchBar extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
      <div className="big-search">
        <input type="text" value="Search for vote, election, minter..." />
        <div className="submit-search"></div>
      </div>
    );
  }
}

