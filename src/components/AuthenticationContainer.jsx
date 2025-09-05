import React, { useState, useEffect, useMemo } from 'react';
import Container from './Container';
import { SubContainer, SubContainerElement } from './Container';
import styles from '../styles/AddGamePage.module.css';

const AuthenticationContainer = ({ 
  actorName, 
  setActorName, 
  actorPassword, 
  setActorPassword,
  players = [],  // Players data passed from parent
  playersLoading = false,
  onUserSelect = null,  // Callback when user is selected (provides user object)
  required = true,
  title = "Authentication",
  titleStyle = "banner"
}) => {
  // Filter and sort auth users (admin and user roles)
  const authUsers = useMemo(() => {
    return players
      .filter(p => p.role === 'admin' || p.role === 'user')
      .sort((a, b) => a.name.localeCompare(b.name));
  }, [players]);

  const handleUserChange = (e) => {
    const userName = e.target.value;
    setActorName(userName);
    
    // Call onUserSelect callback with the selected user object
    if (onUserSelect) {
      const selectedUser = authUsers.find(u => u.name === userName);
      onUserSelect(selectedUser || null);
    }
  };

  return (
    <Container title={title} titleStyle={titleStyle}>
      <SubContainer>

        <SubContainerElement label="Username">
          <select
            name="actor_name"
            className={styles.optionInput}
            value={actorName}
            onChange={handleUserChange}
            required={required}
            disabled={playersLoading}
            autoComplete="username"
          >
            <option value="">Select User</option>
            {authUsers.map((user) => (
              <option key={user.id} value={user.name}>
                {user.name}
              </option>
            ))}
          </select>
        </SubContainerElement>

        {/* Hidden username field for password managers */}
        {actorName && (
          <input 
            type="hidden" 
            name="username" 
            value={actorName}
            autoComplete="username"
          />
        )}

        <SubContainerElement label="Password">
          <input
            type="password"
            name="actor_password"
            className={styles.optionInput}
            value={actorPassword}
            onChange={(e) => setActorPassword(e.target.value)}
            placeholder={actorName ? "Enter password" : "Select user first"}
            disabled={!actorName}
            required={required}
            autoComplete="current-password"
          />
        </SubContainerElement>
      </SubContainer>
    </Container>
  );
};

export default AuthenticationContainer;