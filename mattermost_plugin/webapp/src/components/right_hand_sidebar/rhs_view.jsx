import React from 'react';
import PropTypes from 'prop-types';

import MattarmostPlugin from 'src/mattermost_plugin';

export default class RHSView extends React.PureComponent {
    static propTypes = {
        team: PropTypes.object.isRequired,
    }

    render() {
        return (
            <div style={style.rhs}>
                <MattarmostPlugin/>
            </div>
        );
    }
}

const style = {
    rhs: {
        padding: '10px',
        overflow: 'auto',
    },
};
