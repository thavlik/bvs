import * as React from "react";

import "./../assets/scss/TopSearchBar.scss";

export class TopSearchBar extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
    <div className="search">
        <input type="text" />
        <div className="submit-search"></div>
      </div>
    );
  }
}

