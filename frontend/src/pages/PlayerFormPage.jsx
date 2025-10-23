import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Layout from '../components/Layout';
import Container from '../components/Container';
import { SubContainer, SubContainerElement } from '../components/Container';
import AuthenticationContainer from '../components/AuthenticationContainer';
import LinkButton from '../components/LinkButton';
import { formStyles } from '../styles/formStyles';
import styles from '../styles/GamePage.module.css';
import { playerApi } from '../services/api';

const PlayerFormPage = () => {
  const navigate = useNavigate();
  const { playerId } = useParams();
  const isEditMode = Boolean(playerId);
  const [formData, setFormData] = useState({
    name: '',
    password: '',
    role: 'player',
    actor_name: '',
    actor_password: ''
  });
  const [players, setPlayers] = useState([]);
  const [playersLoading, setPlayersLoading] = useState(true);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [selectedUserRole, setSelectedUserRole] = useState(null);
  const [isInitialSetup, setIsInitialSetup] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch all players for the dropdown
        const data = await playerApi.getAll();
        setPlayers(data);

        // Check if data is empty (no users exist)
        if (data.length === 0) {
          setIsInitialSetup(true);
          // Default to admin role for first user
          setFormData(prev => ({
            ...prev,
            role: 'admin'
          }));
        }

        // If edit mode, fetch player data
        if (isEditMode) {
          const playerData = await playerApi.getById(playerId);
          setFormData(prev => ({
            ...prev,
            name: playerData.name,
            role: playerData.role
          }));
        }
      } catch (err) {
        console.error('Failed to fetch data:', err);
        setError(err.message);
      } finally {
        setPlayersLoading(false);
      }
    };

    fetchData();
  }, [isEditMode, playerId]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;

    // If selecting a user for authentication, track their role
    if (name === 'actor_name') {
      const selectedPlayer = players.find(p => p.name === value);
      setSelectedUserRole(selectedPlayer ? selectedPlayer.role : null);

      // In edit mode, don't change role if editing the authenticated user
      const isEditingSelf = isEditMode && formData.name === value;

      // If the selected user is not an admin and current role is 'user', switch to 'player'
      // Exception: don't change role when editing your own player
      if (selectedPlayer && selectedPlayer.role !== 'admin' && formData.role === 'user' && !isEditingSelf) {
        setFormData(prev => ({
          ...prev,
          actor_name: value,
          role: 'player'
        }));
        return;
      }
    }

    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    // Prepare request data
    const requestData = {
      name: formData.name,
      role: formData.role
    };

    // Include authentication only if not initial setup
    if (!isInitialSetup) {
      requestData.actor_name = formData.actor_name;
      requestData.actor_password = formData.actor_password;
    }

    // Include password for admin and user roles (required for auth roles, not allowed for players)
    // In edit mode, only include password if it's not empty (optional update)
    if (formData.role === 'admin' || formData.role === 'user') {
      if (!isEditMode || formData.password) {
        requestData.password = formData.password;
      }
    }

    try {
      if (isEditMode) {
        await playerApi.update(playerId, requestData);
        alert(`Player "${formData.name}" updated successfully!`);
        navigate(`/players/${playerId}`);
      } else {
        await playerApi.create(requestData);
        alert(`Player "${formData.name}" created successfully!`);
        navigate('/');
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      {!isInitialSetup && (
        <AuthenticationContainer
          actorName={formData.actor_name}
          setActorName={(name) => setFormData(prev => ({ ...prev, actor_name: name }))}
          actorPassword={formData.actor_password}
          setActorPassword={(password) => setFormData(prev => ({ ...prev, actor_password: password }))}
          players={players}
          playersLoading={playersLoading}
          onUserSelect={(user) => {
            setSelectedUserRole(user ? user.role : null);
            // In edit mode, don't change role if editing the authenticated user
            const isEditingSelf = isEditMode && formData.name === user?.name;
            // If the selected user is not an admin and current role is 'user', switch to 'player'
            // Exception: don't change role when editing your own player
            if (user && user.role !== 'admin' && formData.role === 'user' && !isEditingSelf) {
              setFormData(prev => ({ ...prev, role: 'player' }));
            }
          }}
          required={true}
        />
      )}

      <Container title={isInitialSetup ? "Admin Details" : (isEditMode ? "Edit Player" : "New Player Details")} titleStyle="banner">
        <form onSubmit={handleSubmit}>
          <SubContainer>
            <SubContainerElement label={isInitialSetup ? "Admin Name" : "Player Name"}>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleInputChange}
                required
                className={styles.optionInput}
                autoComplete="off"
              />
            </SubContainerElement>

            {!isInitialSetup && (
              <SubContainerElement label="Role">
                <select
                  name="role"
                  value={formData.role}
                  onChange={handleInputChange}
                  className={styles.optionInput}
                  disabled={!formData.actor_name || selectedUserRole === 'user'}
                >
                  <option value="player">Player</option>
                  <option value="user">User</option>
                  {(selectedUserRole === 'admin' || formData.role === 'admin') && (
                    <option value="admin">Admin</option>
                  )}
                </select>
              </SubContainerElement>
            )}

            {(formData.role === 'user' || formData.role === 'admin') && (
              <SubContainerElement label={isEditMode ? "New Password" : "Password"}>
                <input
                  type="password"
                  name="password"
                  value={formData.password}
                  onChange={handleInputChange}
                  required={!isEditMode}
                  className={styles.optionInput}
                  autoComplete="new-password"
                />
              </SubContainerElement>
            )}

            {error && (
              <div style={{
                background: 'rgba(255, 0, 0, 0.1)',
                border: '2px solid rgba(255, 0, 0, 0.5)',
                borderRadius: '10px',
                padding: '1rem',
                color: '#ff0000',
                margin: '1rem 2rem',
                fontSize: '1.5rem'
              }}>
                {error}
              </div>
            )}
          </SubContainer>
        </form>
      </Container>

      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        gap: '2rem',
        margin: '1% auto',
        maxWidth: '900px',
        padding: '0'
      }}>
        <LinkButton
          onClick={() => navigate(isEditMode ? `/players/${playerId}` : '/')}
          style={{ width: 'calc(50% - 1rem)' }}
        >
          Cancel
        </LinkButton>
        <LinkButton
          onClick={handleSubmit}
          disabled={loading}
          style={{ width: 'calc(50% - 1rem)', backgroundColor: '#4CAF50' }}
        >
          {loading ? (isEditMode ? 'Updating...' : 'Creating...') : (isEditMode ? 'Update' : 'Create')}
        </LinkButton>
      </div>
    </Layout>
  );
};

export default PlayerFormPage;
