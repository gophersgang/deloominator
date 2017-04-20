import React, { Component } from 'react';
import VegaLite from 'react-vega-lite';

const spec = {
  "description": "A simple bar chart with embedded data.",
  "mark": "bar",
  "encoding": {
    "x": {"type": "ordinal"},
    "y": {"type": "quantitative"}
  }
};

export default class SimpleBar extends Component {
  render() {
    const data = {
      "values": this.props.values
    };

    // not sure this is a safe assumption
    const [y, x] = Object.keys(this.props.values[0]);

    spec["encoding"]["x"]["field"] = x;
    spec["encoding"]["y"]["field"] = y;

    return (<VegaLite spec={spec} data={data} />);
  }
}
