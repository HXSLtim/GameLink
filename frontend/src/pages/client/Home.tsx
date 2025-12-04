import { Button, Row, Col, Statistic } from 'antd';
import { FireOutlined, SafetyCertificateOutlined, ThunderboltOutlined } from '@ant-design/icons';
import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';

const ClientHome = () => {
    const navigate = useNavigate();

    const containerVariants = {
        hidden: { opacity: 0 },
        visible: {
            opacity: 1,
            transition: {
                staggerChildren: 0.2
            }
        }
    };

    const itemVariants = {
        hidden: { opacity: 0, y: 20 },
        visible: { opacity: 1, y: 0 }
    };

    return (
        <div style={{ minHeight: '100vh', background: 'var(--bg-primary)' }}>
            {/* Hero Section */}
            <section style={{
                height: '80vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: 'var(--bg-primary)',
                position: 'relative',
                overflow: 'hidden'
            }}>
                <div style={{
                    position: 'absolute',
                    top: '20%',
                    left: '10%',
                    width: '300px',
                    height: '300px',
                    background: 'var(--brand-experiment)',
                    filter: 'blur(150px)',
                    opacity: 0.2,
                    borderRadius: '50%'
                }} />
                <div style={{
                    position: 'absolute',
                    bottom: '20%',
                    right: '10%',
                    width: '400px',
                    height: '400px',
                    background: '#eb2f96',
                    filter: 'blur(150px)',
                    opacity: 0.15,
                    borderRadius: '50%'
                }} />

                <motion.div
                    variants={containerVariants}
                    initial="hidden"
                    animate="visible"
                    style={{ textAlign: 'center', zIndex: 1, padding: '0 20px' }}
                >
                    <motion.h1
                        variants={itemVariants}
                        style={{
                            fontSize: '64px',
                            fontWeight: 800,
                            marginBottom: '24px',
                            color: 'var(--text-normal)',
                            lineHeight: 1.2
                        }}
                    >
                        Level Up Your <br />
                        <span style={{ color: 'var(--brand-experiment)' }}>Gaming Experience</span>
                    </motion.h1>
                    <motion.p
                        variants={itemVariants}
                        style={{ fontSize: '20px', color: 'var(--text-muted)', maxWidth: '600px', margin: '0 auto 40px' }}
                    >
                        Find elite companions, join pro squads, and dominate the leaderboard.
                        The ultimate platform for gamers.
                    </motion.p>
                    <motion.div variants={itemVariants}>
                        <Button
                            type="primary"
                            size="large"
                            style={{
                                height: '56px',
                                padding: '0 40px',
                                fontSize: '18px',
                                borderRadius: '28px',
                                backgroundColor: 'var(--brand-experiment)',
                                border: 'none',
                                marginRight: '16px'
                            }}
                            onClick={() => navigate('/companions')}
                        >
                            Find Companion
                        </Button>
                        <Button
                            size="large"
                            ghost
                            style={{
                                height: '56px',
                                padding: '0 40px',
                                fontSize: '18px',
                                borderRadius: '28px',
                                color: 'var(--text-normal)',
                                borderColor: 'var(--text-muted)'
                            }}
                        >
                            Become a Pro
                        </Button>
                    </motion.div>
                </motion.div>
            </section>

            {/* Stats Section */}
            <section style={{ padding: '80px 0', background: 'var(--bg-secondary)' }}>
                <Row justify="center" gutter={[48, 48]}>
                    <Col span={6} style={{ textAlign: 'center' }}>
                        <Statistic
                            title={<span style={{ color: 'var(--text-muted)' }}>Active Players</span>}
                            value={12500}
                            styles={{ content: { color: 'var(--text-normal)', fontSize: '36px', fontWeight: 'bold' } }}
                            suffix="+"
                        />
                    </Col>
                    <Col span={6} style={{ textAlign: 'center' }}>
                        <Statistic
                            title={<span style={{ color: 'var(--text-muted)' }}>Pro Companions</span>}
                            value={850}
                            styles={{ content: { color: 'var(--text-normal)', fontSize: '36px', fontWeight: 'bold' } }}
                            suffix="+"
                        />
                    </Col>
                    <Col span={6} style={{ textAlign: 'center' }}>
                        <Statistic
                            title={<span style={{ color: 'var(--text-muted)' }}>Games Supported</span>}
                            value={45}
                            styles={{ content: { color: 'var(--text-normal)', fontSize: '36px', fontWeight: 'bold' } }}
                            suffix="+"
                        />
                    </Col>
                </Row>
            </section>

            {/* Features Section */}
            <section style={{ padding: '100px 40px', maxWidth: '1200px', margin: '0 auto' }}>
                <div style={{ textAlign: 'center', marginBottom: '80px' }}>
                    <h2 style={{ fontSize: '36px', color: 'var(--text-normal)', marginBottom: '16px' }}>Why Choose GameLink?</h2>
                    <p style={{ fontSize: '16px', color: 'var(--text-muted)' }}>We provide the best ecosystem for gamers to connect and play.</p>
                </div>
                <Row gutter={[32, 32]}>
                    {[
                        {
                            icon: <SafetyCertificateOutlined style={{ fontSize: '32px', color: 'var(--success)' }} />,
                            title: 'Verified Pros',
                            desc: 'All companions undergo strict skill verification and identity checks.'
                        },
                        {
                            icon: <ThunderboltOutlined style={{ fontSize: '32px', color: 'var(--warning)' }} />,
                            title: 'Instant Matching',
                            desc: 'Find a teammate in seconds with our smart matching algorithm.'
                        },
                        {
                            icon: <FireOutlined style={{ fontSize: '32px', color: 'var(--danger)' }} />,
                            title: 'Top Tier Gameplay',
                            desc: 'Play with Challenger/Predator rank players to boost your skills.'
                        }
                    ].map((item, index) => (
                        <Col xs={24} md={8} key={index}>
                            <motion.div
                                whileHover={{ y: -10 }}
                                style={{
                                    background: 'var(--bg-tertiary)',
                                    padding: '40px',
                                    borderRadius: '16px',
                                    height: '100%',
                                    border: '1px solid var(--background-modifier-active)'
                                }}
                            >
                                <div style={{ marginBottom: '24px' }}>{item.icon}</div>
                                <h3 style={{ fontSize: '20px', color: 'var(--text-normal)', marginBottom: '16px' }}>{item.title}</h3>
                                <p style={{ color: 'var(--text-muted)', lineHeight: 1.6 }}>{item.desc}</p>
                            </motion.div>
                        </Col>
                    ))}
                </Row>
            </section>

            {/* Popular Games */}
            <section style={{ padding: '100px 40px', background: 'var(--bg-secondary)' }}>
                <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'end', marginBottom: '40px' }}>
                        <div>
                            <h2 style={{ fontSize: '36px', color: 'var(--text-normal)', marginBottom: '8px' }}>Popular Games</h2>
                            <p style={{ color: 'var(--text-muted)' }}>Trending titles our community loves</p>
                        </div>
                        <Button type="text" style={{ color: 'var(--brand-experiment)' }}>View All Games &rarr;</Button>
                    </div>
                    <Row gutter={[24, 24]}>
                        {['League of Legends', 'Valorant', 'Apex Legends', 'Genshin Impact'].map((game, i) => (
                            <Col xs={24} sm={12} md={6} key={i}>
                                <motion.div
                                    whileHover={{ scale: 1.05 }}
                                    style={{
                                        height: '300px',
                                        background: `linear-gradient(to bottom, rgba(0,0,0,0) 0%, rgba(0,0,0,0.8) 100%), url(https://picsum.photos/400/600?random=${i})`,
                                        backgroundSize: 'cover',
                                        borderRadius: '16px',
                                        position: 'relative',
                                        cursor: 'pointer',
                                        overflow: 'hidden'
                                    }}
                                >
                                    <div style={{ position: 'absolute', bottom: '20px', left: '20px' }}>
                                        <h3 style={{ color: '#fff', margin: 0 }}>{game}</h3>
                                        <span style={{ color: 'rgba(255,255,255,0.7)', fontSize: '12px' }}>1.2k Active Companions</span>
                                    </div>
                                </motion.div>
                            </Col>
                        ))}
                    </Row>
                </div>
            </section>

            {/* CTA */}
            <section style={{ padding: '120px 20px', textAlign: 'center', background: 'var(--bg-tertiary)' }}>
                <h2 style={{ fontSize: '48px', color: 'var(--text-normal)', marginBottom: '24px', fontWeight: 800 }}>Ready to Start?</h2>
                <p style={{ fontSize: '18px', color: 'var(--text-muted)', marginBottom: '40px' }}>Join thousands of gamers and find your perfect duo today.</p>
                <Button
                    type="primary"
                    size="large"
                    style={{
                        height: '64px',
                        padding: '0 60px',
                        fontSize: '20px',
                        borderRadius: '32px',
                        background: 'var(--brand-experiment)',
                        color: '#fff',
                        border: 'none',
                        fontWeight: 'bold'
                    }}
                    onClick={() => navigate('/register')}
                >
                    Get Started Now
                </Button>
            </section>
        </div>
    );
};

export default ClientHome;
