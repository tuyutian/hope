import React, { Component } from "react";
import { Tooltip } from "@shopify/polaris";

class Tooltips extends Component {
  render() {
    return (
      <div>
        <Tooltip width={this.props.width} content={this.props.text}>
          <span className="tooltip_statistical_order" style={{ fontSize: "13px" }}>
            {this.props.title}
          </span>
        </Tooltip>
      </div>
    );
  }
}

export default Tooltips;
