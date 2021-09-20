import * as React from "react";

import "./../assets/scss/TopSearchBar.scss";

export class TopSearchBar extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
    <div className="search">
        <input type="text" value="Search for your vote..." />
        <div className="submit-search"></div>
      </div>
    );
  }
}

