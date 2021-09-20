import * as React from "react";

import SearchIcon from "../assets/img/search.png";

import "./SubmitSearchButton.scss"

export class SubmitSearchButton extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
      <div className="submit-search">
        <img src={SearchIcon} alt="Search" />
      </div>
    );
  }
}

