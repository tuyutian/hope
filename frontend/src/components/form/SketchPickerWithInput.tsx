import { useCallback, useEffect, useState } from "react";
import reactCSS from "reactcss";
import { ColorPicker, HSBAColor, hsbToHex, Popover, rgbToHsb, TextField } from "@shopify/polaris";
import { hexToRgb } from "~/utils/tools";

interface Props {
    defaultColor: string,
    onChange: (color:string)=>void
    popoverStyle?: React.CSSProperties
}

export default function SketchPickerWithInput({ defaultColor, onChange, popoverStyle }: Props) {
    const [display, setDisplay] = useState(false);
    const [color, setColor] = useState("");

    const [hsbColor, setHsbColor] = useState({
        hue: 1,
        brightness: 0,
        saturation: 1
    });
    useEffect(function () {
        setColor(defaultColor || "");
    }, [defaultColor]);
    useEffect(function () {
        setHsbColor(rgbToHsb(hexToRgb(color)))
    }, [color])
    const handleClick = () => {
        setDisplay(!display);
    };

    const handleClose = () => {
        setDisplay(false);
    };

    const handleChange = useCallback((color: HSBAColor) => {
        if (isNaN(color.brightness) || isNaN(color.saturation)) {
            setColor("");
            onChange("")
        } else {
            setColor(hsbToHex(color));
            onChange(hsbToHex(color))
        }
    }, [onChange]);

    const styles = reactCSS({
        "default": {
            color: {
                width: 30,
                height: 30,
                borderRadius: "2px",
                border: "1px solid #999",
                background: color || "url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAwAAAAMCAIAAADZF8uwAAAAGUlEQVQYV2M4gwH+YwCGIasIUwhT25BVBADtzYNYrHvv4gAAAABJRU5ErkJggg==)",
            },
            swatch: {
                padding: "4px",
                background: "#fff",
                borderRadius: "4px",
                border: "1px solid #999",
                boxShadow: "0 0 0 1px rgba(0,0,0,.1)",
                display: "inline-block",
                cursor: "pointer",
            },
            popover: {
                borderRadius: "3px",
                border: 0,
                padding: 0,
                margin: 0,
                display:"flex",
                cursor: "pointer",
                boxShadow: display ? " 0 0 0 2px #fff,0 0 0 4px rgba(69, 143, 255, 1)" : "",
                outline: "none",
                ...popoverStyle,
            },
            cover: {
                position: "fixed",
                top: 0,
                right: 0,
                bottom: 0,
                left: 0,
            },
            colorButton: {
                borderRadius: "3px",
                height: "2rem",
                width: "2rem",
                background: "repeating-conic-gradient(rgba(255, 255, 255, 1) 0 25%,rgba(246, 246, 247, 1) 0 50%) 50%/0.5rem 0.5rem",
                boxShadow: "inset 0 0 0 1px #00000030",
                flexShrink: 0
            },
            colorShow: {
                borderRadius: "inherit",
                boxShadow: "inherit",
                height: "100%",
                width: "100%",
                background: color || "url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAwAAAAMCAIAAADZF8uwAAAAGUlEQVQYV2M4gwH+YwCGIasIUwhT25BVBADtzYNYrHvv4gAAAABJRU5ErkJggg==)",
            }
        },
    });

    return <div >
          <TextField connectedLeft={<Popover active={display} activator={<button onClick={handleClick} type="button" style={styles.popover}>
              <div style={styles.colorButton}>
                  <div style={styles.colorShow} />
              </div>
          </button>
          } onClose={handleClose}>
              <div className="rounded p-2">
                  <ColorPicker color={hsbColor} onChange={(e) => {
                      setHsbColor(e);
                      handleChange(e)
                  }} />
              </div>
          </Popover>} label="" prefix="#" autoComplete="off" value={color.replace('#', '')} onChange={(val: string) => {
              if (val.length === 0) {
                  setColor("");
                  onChange("");
              } else {
                  setColor(`#${val}`);
                  onChange(`#${val}`);
              }
          }} />
      </div>;
}
