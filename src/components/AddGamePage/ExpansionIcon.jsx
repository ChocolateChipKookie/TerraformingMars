import React from "react";
import styles from "../../styles/AddGamePage.module.css";

function ExpansionIcon({
  expansion,
  checked,
  disabled,
  onChange,
  children,
  showText = true,
}) {
  const getExpansionConfig = (expansion) => {
    switch (expansion) {
      case "Base Game":
        return {
          backgroundColor: "transparent",
          symbolColor: "#000",
          svg: null,
        };
      case "Corporate Era":
        return {
          backgroundColor: "#FF0000",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "0deg",
        };
      case "Venus Next":
        return {
          backgroundColor: "#87CEEB",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "180deg",
        };
      case "Prelude":
        return {
          backgroundColor: "#FFC0CB",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "270deg",
        };
      case "Prelude 2":
        return {
          backgroundColor: "#FF69B4",
          symbolColor: "#FFF",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "270deg",
        };
      case "Colonies":
        return {
          backgroundColor: "#808080",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "0deg",
        };
      case "Turmoil":
        return {
          backgroundColor: "#FFA500",
          symbolColor: "#000",
          svg: <polygon points="12,4 20,18 4,18" fill="currentColor" />,
          rotation: "180deg",
        };
      case "Milestones & Awards":
        return {
          backgroundColor: "#FFD700",
          symbolColor: "#000",
          svg: <rect x="6" y="6" width="12" height="12" fill="currentColor" />,
        };
      case "Promo":
        return {
          backgroundColor: "#2F2F2F",
          symbolColor: "#FFF",
          svg: <circle cx="12" cy="12" r="8" fill="currentColor" />,
        };
      default:
        return {
          backgroundColor: "#666",
          symbolColor: "#FFF",
          svg: null,
          fallbackSymbol: "?",
        };
    }
  };

  const config = getExpansionConfig(expansion);
  const isShaded = !checked;

  return (
    <label
      className={
        showText
          ? disabled
            ? styles.expansionLabelExpandedDisabled
            : styles.expansionLabelExpanded
          : disabled
            ? styles.expansionLabelCompactDisabled
            : styles.expansionLabelCompact
      }
    >
      <input
        type="checkbox"
        checked={checked}
        disabled={disabled}
        onChange={onChange}
        className={styles.hiddenCheckbox}
      />
      <div
        className={
          isShaded
            ? styles.expansionIconShaded
            : disabled
              ? styles.expansionIconDisabled
              : styles.expansionIcon
        }
        style={{
          backgroundColor: config.backgroundColor,
          color: config.symbolColor,
        }}
        title={expansion}
      >
        {config.svg ? (
          <svg
            className={styles.expansionIconSvg}
            viewBox="0 0 24 24"
            style={{
              transform: config.rotation
                ? `rotate(${config.rotation})`
                : undefined,
            }}
          >
            {config.svg}
          </svg>
        ) : (
          config.fallbackSymbol || ""
        )}
      </div>
      {showText && <span>{children}</span>}
    </label>
  );
}

export default ExpansionIcon;