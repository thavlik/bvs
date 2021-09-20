import * as React from "react";

import { SubmitSearchButton } from "./SubmitSearchButton";
import "./TopSearchBar.scss";

export class TopSearchBar extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
    <div className="search">
        <input type="text" value="Search for your vote..." />
        <SubmitSearchButton />
      </div>
    );
  }
}

