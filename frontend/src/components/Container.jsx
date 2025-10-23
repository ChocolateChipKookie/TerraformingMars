import React from 'react';
import { formStyles } from '../styles/formStyles';
import styles from '../styles/Container.module.css';

function Container({ children, title, titleStyle = "page-title" }) {
  return (
    <div className={styles.container}>
      {title && titleStyle === "page-title" && (
        <div className={styles.pageTitle}>
          {title}
        </div>
      )}
      {title && titleStyle === "banner" && (
        <div className={styles.bannerTitle}>
          {title}
        </div>
      )}
      {children}
    </div>
  );
}

function SubContainer({ children }) {
  return (
    <div className={styles.subContainer}>
      {children}
    </div>
  );
}

function SubContainerElement({ children, label, input }) {
  return (
    <div className={styles.subContainerElement}>
      {label && <label>{label}</label>}
      {input && (
        <input
          {...input}
          style={{
            ...formStyles.optionInput,
            ...input.style
          }}
        />
      )}
      {children}
    </div>
  );
}

export default Container;
export { SubContainer, SubContainerElement };
