import React, {useState, useEffect} from 'react';

import PropTypes from 'prop-types';

const MattarmostPlugin = ({messages: initialMessages}) => {
    const [messages, setMessages] = useState(initialMessages);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalCount, setTotalCount] = useState(0);
    const [startDate, setStartDate] = useState('');
    const [endDate, setEndDate] = useState('');

    const fetchMessages = async (page, fromDate, toDate) => {
        try {
            let url = `http://localhost:8065/plugins/wfh-plugin/list?page=${page}&per_page=10`;

            // let url = `https://mattermost.vivasoftltd.com/plugins/wfh-plugin/list?page=${page}&per_page=10`;
            if (fromDate) {
                url += `&from_date=${fromDate}`;
            }
            if (toDate) {
                url += `&to_date=${toDate}`;
            }

            const response = await fetch(url);
            const data = await response?.json();

            if (Array.isArray(data?.messages)) {
                setMessages(data?.messages);
                setTotalCount(data?.total_count);
            } else {
                setMessages([]);
                setTotalCount(0);
            }
        } catch (error) {
            setMessages([]);
            setTotalCount(0);
        }
    };

    useEffect(() => {
        fetchMessages(currentPage, startDate, endDate);
    }, [currentPage]);

    const handlePreviousPage = () => {
        if (currentPage > 1) {
            setCurrentPage(currentPage - 1);
        }
    };

    const handleNextPage = () => {
        setCurrentPage(currentPage + 1);
    };

    const handleFilter = () => {
        setCurrentPage(1);

        fetchMessages(1, startDate, endDate);
    };

    return (
        <div>
            <h2 style={{textAlign: 'center'}}>{'Plugin Messages'}</h2>

            <div style={{display: 'flex', justifyContent: 'center', marginBottom: '20px'}}>
                <div style={{marginRight: '10px'}}>
                    <label>{'From Date: '}</label>
                    <input
                        type='date'
                        value={startDate}
                        onChange={(e) => setStartDate(e.target.value)}
                    />
                </div>
                <div style={{marginRight: '10px'}}>
                    <label>{'To Date: '}</label>
                    <input
                        type='date'
                        value={endDate}
                        onChange={(e) => setEndDate(e.target.value)}
                    />
                </div>
                <button
                    onClick={handleFilter}
                    style={{marginLeft: '10px'}}
                >
                    {'Filter'}
                </button>
            </div>

            {messages?.length === 0 ? (
                <p style={{textAlign: 'center'}}>{'No messages found for the selected criteria.'}</p>
            ) : (
                <>
                    <p style={{textAlign: 'center'}}>{`Total Records: ${totalCount}`}</p>

                    <table
                        style={{
                            width: '100%',
                            borderCollapse: 'collapse',
                            marginTop: '20px',
                        }}
                    >
                        <thead>
                            <tr>
                                <th style={style.plugin}>{'Username'}</th>
                                <th style={style.plugin}>{'Team Name'}</th>
                                <th style={style.plugin}>{'Channel Name'}</th>
                                <th style={style.plugin}>{'Message'}</th>
                                <th style={style.plugin}>{'Timestamp'}</th>
                            </tr>
                        </thead>
                        <tbody>
                            {messages?.map((message, index) => (
                                <tr key={index}>
                                    <td style={style.plugin}>{message.username}</td>
                                    <td style={style.plugin}>{message.team_name}</td>
                                    <td style={style.plugin}>{message.channel_name}</td>
                                    <td style={style.plugin}>{message.message}</td>
                                    <td style={style.plugin}>
                                        {new Date(message.timestamp).toLocaleString()}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </>
            )}

            <div style={{display: 'flex', justifyContent: 'center', marginTop: '20px'}}>
                <button
                    onClick={handlePreviousPage}
                    disabled={currentPage === 1}
                >
                    {'Previous'}
                </button>
                <span style={{margin: '0 10px'}}>{`Page ${currentPage}`}</span>
                <button
                    onClick={handleNextPage}
                    disabled={messages?.length < 10}
                >
                    {'Next'}
                </button>
            </div>
        </div>
    );
};

MattarmostPlugin.propTypes = {
    messages: PropTypes.arrayOf(
        PropTypes.shape({
            channel_id: PropTypes.string.isRequired,
            channel_name: PropTypes.string.isRequired,
            team_id: PropTypes.string.isRequired,
            team_name: PropTypes.string.isRequired,
            message: PropTypes.string.isRequired,
            timestamp: PropTypes.number.isRequired,
            user_id: PropTypes.string.isRequired,
            username: PropTypes.string.isRequired,
        }),
    ).isRequired,
};

const style = {
    plugin: {
        border: '1px solid black',
        padding: '8px',
        textAlign: 'center',
    },
};

export default MattarmostPlugin;
