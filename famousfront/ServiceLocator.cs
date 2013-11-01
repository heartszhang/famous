using famousfront.viewmodels;
using System;
using System.Collections.Generic;
using System.Diagnostics.CodeAnalysis;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront
{
    internal class ServiceLocator
    {
        private static MainViewModel _main;

        private static SettingsViewModel _settings;

        /// <summary>
        /// Initializes a new instance of the ViewModelLocator class.
        /// </summary>
        public ServiceLocator()
        {
            CreateMain();
        }

        /// <summary>
        /// Gets the Main property.
        /// </summary>
        internal static MainViewModel Main
        {
            get
            {
                if (_main == null)
                {
                    CreateMain();
                }

                return _main;
            }
        }

        /// <summary>
        /// Gets the Settings property.
        /// </summary>
        internal static SettingsViewModel Settings
        {
            get
            {
                if (_settings == null)
                {
                    CreateSettings();
                }

                return _settings;
            }
        }

        /// <summary>
        /// Gets the Main property.
        /// </summary>
        [SuppressMessage("Microsoft.Performance",
            "CA1822:MarkMembersAsStatic",
            Justification = "This non-static member is needed for data binding purposes.")]
        internal MainViewModel MainViewModel
        {
            get
            {
                return Main;
            }
        }

        /// <summary>
        /// Gets the Settings property.
        /// </summary>
        [SuppressMessage("Microsoft.Performance",
            "CA1822:MarkMembersAsStatic",
            Justification = "This non-static member is needed for data binding purposes.")]
        internal SettingsViewModel SettingsViewModel
        {
            get
            {
                return Settings;
            }
        }

        /// <summary>
        /// Provides a deterministic way to delete the Main property.
        /// </summary>
        internal static void ClearMain()
        {
            _main = null;
        }

        /// <summary>
        /// Provides a deterministic way to delete the Settings property.
        /// </summary>
        internal static void ClearSettings()
        {
            _settings = null;
        }

        /// <summary>
        /// Provides a deterministic way to create the Main property.
        /// </summary>
        internal static void CreateMain()
        {
            if (_main == null)
            {
                _main = new MainViewModel();
            }
        }

        /// <summary>
        /// Provides a deterministic way to create the Settings property.
        /// </summary>
        internal static void CreateSettings()
        {
            if (_settings == null)
            {
                _settings = new SettingsViewModel();
            }
        }
    }
}
