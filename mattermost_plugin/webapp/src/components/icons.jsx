import React, {useState, useEffect} from 'react';
import {Client4} from 'mattermost-redux/client';

Client4.setUrl('http://localhost:8065');

// Client4.setUrl('https://mattermost.vivasoftltd.com');

export const ChannelHeaderButtonIcon = () => {
    const [isSystemAdmin, setIsSystemAdmin] = useState<boolean | null>(null);

    useEffect(() => {
        const fetchUserRoles = async () => {
            try {
                const user = await Client4.getMe();
                const roles = user.roles.split(' ');
                setIsSystemAdmin(roles.includes('system_admin'));
            } catch (error) {
                setIsSystemAdmin(false);
            }
        };

        fetchUserRoles();
    }, []);

    if (isSystemAdmin === null) {
        return null;
    }

    if (isSystemAdmin) {
        return (
            <div>
                <i
                    className='icon fa fa-plug'
                    style={{fontSize: '15px', position: 'relative', top: '-1px'}}
                />
            </div>

        );
    }

    return null;
};
