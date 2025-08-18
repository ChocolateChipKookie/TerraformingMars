import React from 'react';

function Layout({ children }) {
  return (
    <div style={{
      backgroundColor: '#021631',
      background: 'url(/resources/mars.png) bottom no-repeat, linear-gradient(#010b19, #010b19, #03214a, #85c7f1)',
      backgroundSize: '100% auto',
      backgroundAttachment: 'fixed',
      fontFamily: 'tm-default',
      textTransform: 'uppercase',
      paddingLeft: '5%',
      paddingRight: '5%',
      minHeight: '100vh',
      position: 'relative',
      overflowY: 'auto'
    }}>
      {children}
    </div>
  );
}

export default Layout;
