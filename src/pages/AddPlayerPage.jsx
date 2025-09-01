import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import Container from '../components/Container';
import { SubContainer, SubContainerElement } from '../components/Container';
import LinkButton from '../components/LinkButton';
import { formStyles } from '../styles/formStyles';
import styles from '../styles/AddGamePage.module.css';

const AddPlayerPage = () => {
  const navigate = useNavigate();
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

  useEffect(() => {
    // Fetch all players for the dropdown
    const fetchPlayers = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/players');
        if (response.ok) {
          const data = await response.json();
          setPlayers(data);
        }
      } catch (err) {
        console.error('Failed to fetch players:', err);
      } finally {
        setPlayersLoading(false);
      }
    };

    fetchPlayers();
  }, []);

  const handleInputChange = (e) => {
    const { name, value } = e.target;

    // If selecting a user for authentication, track their role
    if (name === 'actor_name') {
      const selectedPlayer = players.find(p => p.name === value);
      setSelectedUserRole(selectedPlayer ? selectedPlayer.role : null);

      // If the selected user is not an admin and current role is 'user', switch to 'player'
      if (selectedPlayer && selectedPlayer.role !== 'admin' && formData.role === 'user') {
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
      role: formData.role,
      actor_name: formData.actor_name,
      actor_password: formData.actor_password
    };

    // Only include password for users (required for users, not allowed for players)
    if (formData.role === 'user') {
      requestData.password = formData.password;
    }

    try {
      const response = await fetch('http://localhost:8080/api/players', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData)
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to create player');
      }

      alert(`Player "${formData.name}" created successfully!`);
      navigate('/');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <Container title="Add New Player" titleStyle="banner">
        <form onSubmit={handleSubmit}>
          <SubContainer>
            <div style={formStyles.subcontainerBox}>
              <SubContainer>
                <div style={{ fontSize: '2rem', marginBottom: '1rem' }}>
                  Authentication
                </div>

                <SubContainerElement label="Username">
                  <select
                    name="actor_name"
                    value={formData.actor_name}
                    onChange={handleInputChange}
                    required
                    className={styles.optionInput}
                    disabled={playersLoading}
                    autoComplete="username"
                  >
                    <option value="">Select User</option>
                    {players
                      .filter(player => player.role === 'user' || player.role === 'admin')
                      .map(player => (
                        <option key={player.id} value={player.name}>
                          {player.name}
                        </option>
                      ))}
                  </select>
                </SubContainerElement>

                {/* Hidden username field for password managers */}
                {formData.actor_name && (
                  <input 
                    type="hidden" 
                    name="username" 
                    value={formData.actor_name}
                    autoComplete="username"
                  />
                )}

                <SubContainerElement label="Password">
                  <input
                    type="password"
                    name="actor_password"
                    value={formData.actor_password}
                    onChange={handleInputChange}
                    required
                    className={styles.optionInput}
                    autoComplete="current-password"
                  />
                </SubContainerElement>
              </SubContainer>
            </div>

            <div style={formStyles.subcontainerBox}>
              <SubContainer>
                <div style={{ fontSize: '2rem', marginBottom: '1rem' }}>
                  New Player Details
                </div>

                <SubContainerElement label="Player Name">
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
                  </select>
                </SubContainerElement>

                {formData.role === 'user' && (
                  <SubContainerElement label="Password">
                    <input
                      type="password"
                      name="password"
                      value={formData.password}
                      onChange={handleInputChange}
                      required
                      className={styles.optionInput}
                      autoComplete="new-password"
                    />
                  </SubContainerElement>
                )}
              </SubContainer>
            </div>

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

            <div style={{
              display: 'flex',
              gap: '2rem',
              justifyContent: 'center',
              padding: '2rem'
            }}>
              <LinkButton
                onClick={() => navigate('/')}
                style={{ width: '200px' }}
              >
                Cancel
              </LinkButton>
              <LinkButton
                onClick={handleSubmit}
                disabled={loading}
                style={{ width: '200px' }}
              >
                {loading ? 'Creating...' : 'Create'}
              </LinkButton>
            </div>
          </SubContainer>
        </form>
      </Container>
    </Layout>
  );
};

export default AddPlayerPage;
