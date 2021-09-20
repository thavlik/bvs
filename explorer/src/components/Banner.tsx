import * as React from "react";
import { TopSearchBar } from "./TopSearchBar";
import "./../assets/scss/Banner.scss";

export class Banner extends React.Component<Record<string, unknown>, undefined> {
  public render() {
    return (
      <div className="banner">
        <div className="banner-content">
          Blockchain Voting Systems
          <div className="nav">
            <div className="nav-item"><a href="#">Elections</a></div>
            <div className="nav-item"><a href="#">Minters</a></div>
            <div className="nav-item"><a href="#">Votes</a></div>
          </div>
          <TopSearchBar />
          <div className="nav">
            <div className="nav-item"><a href="#">FAQ</a></div>
            <div className="nav-item"><a href="#">Support</a></div>
          </div>
        </div>
      </div>
    );
  }
}

