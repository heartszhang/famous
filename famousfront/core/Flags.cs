using GalaSoft.MvvmLight.Messaging;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Runtime.Serialization;
using System.Text;
using System.Threading.Tasks;
using System.Windows;

namespace famousfront.core
{
    [DataContract]
    class Flags
    {
        public Flags()
        {
            Backend = "127.0.0.1:8002";
            KaPeriod = 100;  // milliseconds
            ContentMarginHeight = 16.0;
            FeedSourceWidth = 360.0;
            ContentMarginMinWidth = 180.0;
            TinyFontSize = 12.0;
            SmallFontSize = 13.333;
            NormalFontSize = 16.0;
            LargeFontSize = 17.0;
            HugeFontSize = 18.666;

            DefaultFont = new System.Windows.Media.FontFamily("Microsoft YaHei UI, Segoe UI");
            SemiBoldFont = new System.Windows.Media.FontFamily("Microsoft YaHei UI, Segoe UI Semibold");
            LightFont = new System.Windows.Media.FontFamily("Microsoft YaHei UI Light, Segoe UI Light");

            Theme = Elysium.Theme.Light;
            ThemeAccent = Elysium.AccentBrushes.Blue.Color;
            ThemeContrast = System.Windows.Media.Colors.White;
            MarginTopDown = new Thickness(0, 8, 0, 8);
            MarginAllBounds = new Thickness(4);
            ImageMaxHeight = 640.0;
            ImageMaxWidth = 320.0;
            VideoElementHeight = 331.0;
            BodyMaxWidth = 800.0;
            SearchBoxMargin = new Thickness(160.0, 12, 160, 12);
            FeedSourceFindEntriesMaxWidth = 800;
            ProgressRingLargeSize = 80.0;
        }
        [DataMember(Name="backend")]
        public string Backend
        {
            get;
            set;
        }
        [DataMember(Name = "ka_period")]
        public int KaPeriod
        {
            get;
            set;
        }
        [DataMember(Name = "contentmargin_height")]
        public double ContentMarginHeight { get; set; }

        [DataMember(Name = "feedsource_width")]
        public double FeedSourceWidth { get; set; }

        [DataMember(Name = "contentmargin_minwidth")]
        public double ContentMarginMinWidth { get; set; }

        [DataMember(Name = "tiny_fontsize")]
        public double TinyFontSize { get; set; }
        [DataMember(Name = "small_fontsize")]
        public double SmallFontSize { get; set; }
        [DataMember(Name = "normal_fontsize")]
        public double NormalFontSize { get; set; }
        [DataMember(Name = "large_fontsize")]
        public double LargeFontSize { get; set; }
        [DataMember(Name = "huge_fontsize")]
        public double HugeFontSize { get; set; }

        [DataMember(Name = "default_font")]
        public System.Windows.Media.FontFamily DefaultFont { get; set; }
        [DataMember(Name = "light_font")]
        public System.Windows.Media.FontFamily LightFont { get; set; }
        [DataMember(Name = "semibold_font")]
        public System.Windows.Media.FontFamily SemiBoldFont { get; set; }

        [DataMember(Name = "theme")]
        public Elysium.Theme Theme { get; set; }
        [DataMember(Name = "theme_accent")]
        public System.Windows.Media.Color ThemeAccent { get; set; }
        [DataMember(Name = "theme_contrast")]
        public System.Windows.Media.Color ThemeContrast { get; set; }

        [DataMember(Name = "margin_topdown")]
        public Thickness MarginTopDown { get; set; }
        [DataMember(Name = "margin_allbounds")]
        public Thickness MarginAllBounds { get; set; }
        [DataMember(Name = "image_maxheight")]
        public double ImageMaxHeight{get;set;}

        [DataMember(Name = "image_maxwidth")]
        public double ImageMaxWidth { get; set; }

        [DataMember(Name = "body_maxwidth")]
        public double BodyMaxWidth { get; set; }

        [DataMember(Name = "video_elementheight")]
        public double VideoElementHeight { get; set; }

        [DataMember(Name = "searchbox_margin")]
        public Thickness SearchBoxMargin{ get; set; }

        [DataMember(Name = "feedsource_findentries_maxwidth")]
        public double FeedSourceFindEntriesMaxWidth { get; set; }

        [DataMember(Name = "progressring_largesize")]
        public double ProgressRingLargeSize { get; set; }
    }

}
